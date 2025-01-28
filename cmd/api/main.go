package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kinzaz/social/internal/auth"
	"github.com/kinzaz/social/internal/db"
	"github.com/kinzaz/social/internal/env"
	"github.com/kinzaz/social/internal/mailer"
	"github.com/kinzaz/social/internal/ratelimiter"
	"github.com/kinzaz/social/internal/store"
	"github.com/kinzaz/social/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// BasePath /v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	if err := env.LoadDotEnv(); err != nil {
		log.Fatal("Error while loading file .env: ", err)
	}

	// Config
	cfg := Config{
		addr:        env.GetString("ADDR", "8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			dsn:          env.GetString("DSN", "postgres://admin:pass@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailGun: mailGunConfig{
				apiKey: env.GetString("MAILGUN_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("AUTH_BASIC_USER", "admin"),
				password: env.GetString("AUTH_BASIC_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", ""),
				exp:    time.Hour * 24 * 3,
				issuer: "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUEST_COUNT", 20),
			TimeFrame:            time.Second * 3,
			Enabled:              env.GetBool("RATELIMITER_ENABLED", true),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindow(cfg.rateLimiter.RequestsPerTimeFrame, cfg.rateLimiter.TimeFrame)

	// store
	store := store.NewStorage(db)
	// cache store
	cacheStorage := cache.NewRedisStorage(rdb)

	// mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	// mailer := mailer.NewMailGun(cfg.mail.mailGun.apiKey, cfg.mail.fromEmail)
	mailer, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// auth
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
