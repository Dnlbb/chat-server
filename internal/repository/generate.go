package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ./repointerface.StorageInterface -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ./repointerface.AuthInterface -o ./mocks/ -s "_minimock.go"
