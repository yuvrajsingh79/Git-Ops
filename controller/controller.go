package controller

import (
	"almabase/Git-Ops/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//GetGitDetails functions returns the top most popular repositories of an organisation based on number of forks
func GetGitDetails(w http.ResponseWriter, r *http.Request) {
	//here we are usiing github access token for authenticating the api request
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "No Authorization token found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	org := r.URL.Query().Get("Org")
	numRepos, err := strconv.Atoi(r.URL.Query().Get("n"))
	numCommittees, err := strconv.Atoi(r.URL.Query().Get("m"))
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenString},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	if client == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	//fetching the top n most popular list of repositories of a given organisation
	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)
	if repos == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := err.(*github.AcceptedError); ok {
		log.Println("scheduled on GitHub side")
	}
	sort.SliceStable(repos, func(i, j int) bool {
		return repos[i].GetForksCount() > repos[j].GetForksCount()
	})
	var result model.Git
	reposit := []*model.Repo{}
	for i, repo := range repos {
		if i >= numRepos {
			break
		}
		fmt.Printf("RepoName : %s , Forks : %d \n\n", repo.GetName(), repo.GetForksCount())
		rep := new(model.Repo)
		rep.RepoName = repo.GetName()
		rep.Forks = repo.GetForksCount()
		opts := &github.ListContributorsOptions{
			Anon:        "false",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		contributors, _, err := client.Repositories.ListContributors(context.Background(), org, repo.GetName(), opts)
		if _, ok := err.(*github.AcceptedError); ok {
			log.Println("Repositories.ListContributors returned error: ", err)
		}
		res := []*model.Committee{}
		for i, contributor := range contributors {
			if i >= numCommittees {
				break
			}
			fmt.Println("committee : ", contributor.GetLogin(), " -->> ", contributor.GetContributions())

			com := new(model.Committee)
			com.Name = contributor.GetLogin()
			com.Commits = contributor.GetContributions()
			res = append(res, com)
		}
		fmt.Println("--------------********-------------")
		rep.Committee = res
		reposit = append(reposit, rep)
	}
	result.Repo = reposit
	json.NewEncoder(w).Encode(&reposit)
	return
}
