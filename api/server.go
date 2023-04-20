package api

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gb-detection/recognition"
	"gb-detection/store"
	"sync"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
		AuthKey string `yaml:"auth_key"`
	} `yaml:"server"`
	Databases struct {
		Staffs      string `yaml:"staffs"`
		Images      string `yaml:"images"`
		Thirdpartys string `yaml:"thirdpartys"`
		Timerecords string `yaml:"timerecords"`
	} `yaml:"databases"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("`%s` is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}

func ServerRun() {
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	staffDb := prepareStaffDB(cfg.Databases.Staffs)
	defer staffDb.SaveToFile()
	imageDb := prepareImageDB(cfg.Databases.Images)
	defer imageDb.SaveToFile()
	timerecordDb := prepareTimeRecordDB(cfg.Databases.Timerecords)
	defer timerecordDb.SaveToFile()
	thirdpartyDb := prepareThirdpartyDB(cfg.Databases.Thirdpartys)
	defer thirdpartyDb.SaveToFile()
	recognizer, err := recognition.NewRecognizer()
	if err != nil {
		log.Fatal(err)
	}
	defer recognizer.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	authKey := cfg.Server.AuthKey
	mKey := sync.Mutex{}
	router.Use(func(ctx *gin.Context) {
		key := ctx.Request.Header.Get("Authorization")
		mKey.Lock()
		check := false
		if key != authKey {
			check = true
		}
		mKey.Unlock()
		if check {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		ctx.Next()
	})
	ImageGroup(router, imageDb, recognizer, url)
	TimeRecordGroup(router, timerecordDb, url)
	ThirdpartyGroup(router, timerecordDb, thirdpartyDb, url)
	StaffGroup(router, staffDb, imageDb, recognizer, url)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	log.Printf("Server is running on %s", url)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout.Read * time.Second,
		WriteTimeout: cfg.Server.Timeout.Write * time.Second,
		IdleTimeout:  cfg.Server.Timeout.Idle * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}

	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func prepareStaffDB(path string) *store.StaffDBmem {
	db := store.NewStaffDBmem()
	db.LoadFromFile(path)
	return db
}

func prepareImageDB(path string) *store.ImageDBmem {
	db := store.NewImageDBmem()
	db.LoadFromFile(path)
	return db
}

func prepareTimeRecordDB(path string) *store.TimeRecordDBMem {
	db := store.NewTimeRecordDBMem()
	db.LoadFromFile(path)
	return db
}

func prepareThirdpartyDB(path string) *store.ThirdpartyDBMem {
	db := store.NewThirdpartyDBMem()
	db.LoadFromFile(path)
	return db
}
