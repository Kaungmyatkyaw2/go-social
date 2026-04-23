package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/Kaungmyatkyaw2/go-social/internal/store"
)

var usernames = []string{
	"zenith", "vortex", "pixel", "echo", "flint",
	"marrow", "cipher", "glyph", "juno", "atlas",
	"river", "quill", "fable", "onyx", "nova",
	"orbit", "haven", "koda", "lyric", "sage",
	"mobi", "flux", "roam", "veda", "rift",
	"gaze", "bolt", "aura", "neon", "drift",
	"lark", "prism", "haze", "tide", "jinx",
	"solo", "apex", "moth", "fern", "vibe",
	"axis", "bold", "clay", "dusk", "opal",
	"gulf", "iron", "jade", "kite", "loom",
}

var titles = []string{
	"Vibes", "Golden", "Escape", "Mood", "Unfiltered",
	"Wanderlust", "Glow", "Reflections", "Chasing", "Currents",
	"Essentials", "Afterglow", "Fragments", "Serenity", "Stardust",
	"Legendary", "Wild", "Breathe", "Obsessed", "Radiant",
}

var contents = []string{
	"The best views come after the hardest climbs. Grateful for this moment.",
	"Sometimes you just need to unplug and enjoy the silence.",
	"Focusing on the process rather than the outcome today.",
	"Chasing sunsets and better versions of myself.",
	"The city looks different when you stop to actually look at it.",
	"Finding magic in the mundane parts of the daily routine.",
	"Less perfection, more authenticity. That is the goal.",
	"Current mood: coffee in hand and a clear mind.",
	"Reflecting on everything that brought me to this exact spot.",
	"Small steps still move you forward. Keep going.",
	"There is beauty in the things we usually overlook.",
	"New day, new energy, and a whole lot of gratitude.",
	"Just another chapter in this wild, beautiful journey.",
	"Creating a life I do not need a vacation from.",
	"The stars only shine when it is dark enough to see them.",
	"Quiet mornings are the best kind of therapy.",
	"Everything you need is already within you. Believe it.",
	"Surround yourself with people who feel like sunshine.",
	"Trusting the timing of my life a little more every day.",
	"Keep it simple. The rest will fall into place.",
}

var tags = []string{
	"lifestyle", "aesthetic", "vlog", "mindset", "travel",
	"minimal", "daily", "explore", "vibes", "growth",
	"creativity", "nature", "mood", "wellness", "journal",
	"hustle", "moment", "horizon", "focus", "dreamy",
}

var commentContents = []string{
	"Love the vibe of this post!",
	"So true, thanks for sharing.",
	"Where was this taken? Looks amazing.",
	"Needed to hear this today. 💯",
	"This is exactly what I was thinking.",
	"Pure aesthetic. Love it.",
	"Great perspective!",
	"Adding this to my vision board.",
	"Can't wait to see more from you.",
	"Absolutely stunning.",
	"Spot on! Keep it up.",
	"This deserves more likes honestly.",
	"Iconic.",
	"Such a mood right now.",
	"Checking out your profile, great content!",
	"Wow, just wow.",
	"Thanks for the inspiration!",
	"Keep shining! ✨",
	"Totally agree with this.",
	"Simple yet so powerful.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)

	for _, post := range comments {
		if err := store.Comments.Create(ctx, post); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)

	for i := 0; i < num; i++ {

		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)

		users[i] = &store.User{
			Username: username,
			Email:    fmt.Sprintf("%s@example.com", username),
			Password: fmt.Sprintf("%s#123123", username),
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
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}

	}

	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {

		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			UserID:  user.ID,
			Content: commentContents[rand.Intn(len(commentContents))],
			PostID:  post.ID,
		}

	}

	return comments

}
