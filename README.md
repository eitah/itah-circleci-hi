# Itah CircleCI Go

This is a sample go project that I can use for a pipeline project in
circleci.

Tips for eli from spinning this up before

* make sure your null.zip for the lambda is created - see instructions in the terraform
* to install a pakage in go use `go get github.com/sirupsen/logrus`
* rand.Seed(time.Now().UnixNano()) before calling rand.Int63n is overkill lol!