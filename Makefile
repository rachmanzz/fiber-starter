.PHONY: init
init:
	@read -p "Enter New Module Name (e.g., github.com/user/repo): " modname; \
	find . -type f -not -path '*/.*' -exec sed -i "s|github.com/rachmanzz/fiber-starter|$$modname|g" {} +; \
	go mod edit -module $$modname; \
	go mod tidy; \
	echo "Project initialized with module: $$modname"