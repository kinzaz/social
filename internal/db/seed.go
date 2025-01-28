package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/kinzaz/social/internal/store"
)

var usernames = []string{
	"Alex", "Bob", "Charlie", "David", "Emily", "Frank", "Grace", "Henry", "Isabella", "Jack",
	"Kate", "Leo", "Mia", "Noah", "Olivia", "Peter", "Quinn", "Rachel", "Samuel", "Tessa",
	"Ulysses", "Victor", "Wendy", "Xavier", "Yasmin", "Zachary", "Albert", "Beatrice", "Cameron",
	"Daniela", "Ethan", "Felix", "Georgia", "Hugo", "Imogen", "Jordan", "Kyle", "Lucas", "Madeline",
	"Nathan", "Oscar", "Penelope", "Quentin", "Rebecca", "Simon", "Thomas", "Uma", "Vincent", "Willow",
	"Xena", "York", "Zoe",
}

var titles = []string{
	"How to Choose the Perfect Gift?",
	"Journey to Unexplored Places",
	"Top 10 Books You Should Read",
	"Secrets of a Successful Career",
	"Healthy Lifestyle: Myths and Reality",
	"Tips for Taking Care of Houseplants",
	"Modern Technologies in Everyday Life",
	"How to Organize Your Time Effectively?",
	"Culinary Masterpieces for Every Day",
	"The History of Great Inventions",
	"Fascinating Facts About Animals",
	"Best Movies of the Year",
	"Guide to Popular Tourist Destinations",
	"Tips for Improving Sleep Quality",
	"How to Deal with Stress?",
	"New Trends in Fashion",
	"Entertaining Science Experiments",
	"How to Learn Guitar Playing?",
	"Useful Habits for Good Health",
	"Best Apps for Smartphones",
}

var contents = []string{
	"Discover tips and tricks to find the best gifts for any occasion.",
	"Explore unknown places and create unforgettable memories.",
	"Check out our list of must-read books that will broaden your horizons.",
	"Learn the secrets to building a successful and fulfilling career.",
	"Debunk common myths about leading a healthy lifestyle and uncover the truth.",
	"Get expert advice on keeping your houseplants happy and thriving.",
	"Find out how modern technologies are changing the way we live and work.",
	"Master the art of effective time management with these practical tips.",
	"Elevate your everyday meals with easy yet impressive culinary creations.",
	"Delve into the fascinating history behind groundbreaking inventions.",
	"Discover intriguing facts about the animal kingdom that will amaze you.",
	"Stay up-to-date with the best movies released this year.",
	"Plan your next adventure with our comprehensive guide to popular tourist spots.",
	"Improve your sleep quality with these science-backed tips and tricks.",
	"Manage stress effectively with proven techniques and strategies.",
	"Stay ahead of the curve with the latest trends in fashion.",
	"Uncover entertaining scientific experiments you can do at home.",
	"Start your musical journey with our beginner's guide to playing guitar.",
	"Develop healthy habits that will enhance your overall well-being.",
	"Discover the best smartphone apps to simplify and enrich your life.",
}

var tags = []string{
	"gifts", "travel", "books", "career", "health", "houseplants", "technology", "productivity", "cooking", "history", "animals", "movies", "tourism", "sleep", "stress", "fashion", "science", "music", "habits", "apps",
}

var comments = []string{
	"This is such a helpful guide! Thanks for sharing!",
	"I can't wait to explore these places! So inspiring!",
	"Just added a few of these books to my reading list. Excited to dive in!",
	"Great advice! I'll definitely implement some of these career tips.",
	"Loved this article! It's refreshing to see these myths debunked.",
	"My plants are finally thriving thanks to your tips. Thank you!",
	"Technology has indeed changed our lives in so many ways. Interesting read!",
	"Time management has always been a struggle for me. This was really insightful.",
	"Your recipes look amazing! Can't wait to try them out.",
	"Fascinating to learn about these inventions. Great job!",
	"So many interesting facts about animals! Keep up the good work!",
	"Look forward to checking out these movies. They sound fantastic!",
	"Bookmarked this guide for my next vacation planning session. Thanks!",
	"Sleep has improved since following your tips. Highly recommend this article!",
	"Much needed stress relief tips. Appreciate the suggestions!",
	"Love the new fashion trends! Will definitely incorporate some into my wardrobe.",
	"Fun experiments! My kids will love trying these out.",
	"Finally picked up the guitar again. Your guide is super helpful!",
	"These healthy habits are game-changers. Feeling better already!",
	"Downloaded a couple of these apps. They're really useful!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(titles))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {

		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
