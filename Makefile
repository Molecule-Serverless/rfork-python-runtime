.PHONY: FORCE
base-image: FORCE
	docker build -t python-base-image python-base-image

base-fs: FORCE
	sh python-base-image/make-base-fs.sh
