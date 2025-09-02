package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MeYo0o/blog_aggregator/internal/config"
	"github.com/MeYo0o/blog_aggregator/internal/database"
	"github.com/MeYo0o/blog_aggregator/internal/rss"
	st "github.com/MeYo0o/blog_aggregator/internal/state"
	"github.com/google/uuid"
)

func HandlerLogin(s *st.State, cmd Command) error {
	var err error
	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: login
		// Args[2] is the username
		loginUsername := cmd.Args[2]
		_, err = s.DB.GetUser(context.Background(), loginUsername)
		if err != nil {
			return errors.New("user doesn't exist in DB")
		} else {
			s.Cfg.CurrentUsername = loginUsername
			config.SetUser(s.Cfg.CurrentUsername)
		}
	default:
		return errors.New("you need to pass the username only after login")
	}

	fmt.Printf("user %s has been set!\n", s.Cfg.CurrentUsername)

	return nil
}

func HandlerRegister(s *st.State, cmd Command) error {
	var user database.User
	var err error

	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: register
		// Args[2] is the username to be stored inside db
		registrationName := cmd.Args[2]
		user, err = s.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      registrationName,
		})
		if err != nil {
			return errors.New("name already exists")
		} else {
			s.Cfg.CurrentUsername = registrationName
			config.SetUser(registrationName)
		}
	default:
		return errors.New("you need to pass the username only after register")
	}

	fmt.Printf("user %s is stored in DB!\n", s.Cfg.CurrentUsername)
	fmt.Printf("User: %v\n", user)

	return nil
}

func HandleResetUsers(s *st.State, cmd Command) error {
	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: reset
		if err := s.DB.ResetUsers(context.Background()); err != nil {
			return fmt.Errorf("reset Users Failed: %w", err)
		}
	default:
		return errors.New("you don't need any arguments, just the reset command will do")
	}

	fmt.Println("All users have been deleted successfully!")

	return nil
}

func HandleGetUsers(s *st.State, cmd Command) error {
	var users []database.User
	var err error

	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: users
		users, err = s.DB.GetUsers(context.Background())
		if err != nil {
			return errors.New("couldn't retrieve users from DB")
		}

	default:
		return errors.New("you don't need any arguments, just the users command will do")
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}

func HandleAgg(s *st.State, cmd Command) error {
	var rssFeed *rss.RSSFeed
	var err error

	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: agg
		url := "https://www.wagslane.dev/index.xml"
		rssFeed, err = rss.FetchFeed(context.Background(), url)
		if err != nil {
			return fmt.Errorf("error when fetching RSS: %w", err)
		}

	default:
		return errors.New("you don't need any arguments, just the agg command will do")
	}

	fmt.Println(rssFeed)

	return nil
}
func HandleAddFeed(s *st.State, cmd Command, user database.User) error {
	var feedName, feedUrl string

	switch len(cmd.Args) {
	case 4:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: addfeed
		// Args[2] is the command name, i.e: Feed Name
		// Args[3] is the command name, i.e: Feed Url
		feedName = cmd.Args[2]
		feedUrl = cmd.Args[3]

		// add the feed mapped to the user
		feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      feedName,
			Url:       feedUrl,
			UserID:    user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("couldn't create a feed for the current user: %w", err)
		}

		// follow that feed
		_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't follow the feed that was just created for the current user: %w", err)
		}

	default:
		return errors.New("you need to provide 2 more arguments: feedName & feedUrl")
	}

	return nil
}

func HandleGetFeeds(s *st.State, cmd Command) error {
	var feeds []database.Feed
	var err error

	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: users
		feeds, err = s.DB.GetFeeds(context.Background())
		if err != nil {
			return errors.New("couldn't retrieve feeds from DB")
		}

	default:
		return errors.New("you don't need any arguments, just the feeds command will do")
	}

	for _, feed := range feeds {

		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get a feed's attached username: %w", err)
		}

		fmt.Printf("- '%s'\n", feed.Name)
		fmt.Printf("- \"%s\"\n", feed.Url)
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}

func HandleFollowFeed(s *st.State, cmd Command, user database.User) error {
	var feedsFollow []database.FeedFollow

	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: follow feed
		// Args[2] is the command name, i.e: feedUrl
		feedUrl := cmd.Args[2]

		currentFeed, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)
		if err != nil {
			return fmt.Errorf("couldn't find a feed with that Url to follow: %w", err)
		}

		feedsFollow, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    currentFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't follow the provided feedUrl: %w", err)
		}

	default:
		return errors.New("you need to provide 1 more argument: feedUrl")
	}

	for _, feedFollow := range feedsFollow {
		user, err := s.DB.GetUserByID(context.Background(), feedFollow.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get feedFollow linked User: %w", err)
		}
		feed, err := s.DB.GetFeedByID(context.Background(), feedFollow.FeedID)
		if err != nil {
			return fmt.Errorf("couldn't get feedFollow linked Feed: %w", err)
		}

		fmt.Println(user.Name)
		fmt.Println(feed.Name)
	}

	return nil
}

func HandleFollowing(s *st.State, cmd Command, user database.User) error {
	var feedsFollow []database.FeedFollow
	var err error

	switch len(cmd.Args) {
	case 2:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: following
		feedsFollow, err = s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
		if err != nil {
			return fmt.Errorf("couldn't get feed follows for the current user: %w", err)
		}

	default:
		return errors.New("you don't need any arguments, just the following command will do")
	}

	for _, feedFollow := range feedsFollow {
		user, err := s.DB.GetUserByID(context.Background(), feedFollow.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get feedFollow linked User: %w", err)
		}
		feed, err := s.DB.GetFeedByID(context.Background(), feedFollow.FeedID)
		if err != nil {
			return fmt.Errorf("couldn't get feedFollow linked Feed: %w", err)
		}

		fmt.Println(user.Name)
		fmt.Println(feed.Name)
	}

	return nil
}

func HandleUnfollow(s *st.State, cmd Command, user database.User) error {
	var feedUrl string

	switch len(cmd.Args) {
	case 3:
		// Args[0] is the program name, we don't need that but it exists no matter what.
		// Args[1] is the command name, i.e: unfollow
		// Args[2] is the command name, i.e: feedUrl
		feedUrl := cmd.Args[2]
		// get the feed
		feed, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)
		if err != nil {
			return fmt.Errorf("couldn't get the feed with the provided url:%s error:%w", feedUrl, err)
		}

		err = s.DB.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't unfollow the feed with the provided url:%s error:%w", feedUrl, err)
		}

	default:
		return errors.New("you don't need any arguments, just the following command will do")
	}

	fmt.Printf("Feed Url: %s has been unfollowed!\n", feedUrl)

	return nil
}
