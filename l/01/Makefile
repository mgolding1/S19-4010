
all: Lect-01.html README.html

README.html: README.md
	markdown-cli --input=./README.md --output=README.html
	echo cat ../../md.css README.html >/tmp/README.html
	cat ./css/pre ./css/markdown.css ./css/post ../../md.css ./css/hpre README.html ./css/hpost >/tmp/README.html
	mv /tmp/README.html ./README.html

Lect-01.html: Lect-01.md
	markdown-cli --input=./Lect-01.md --output=Lect-01.html
	echo cat ../../md.css Lect-01.html >/tmp/Lect-01.html
	cat ./css/pre ./css/markdown.css ./css/post ../../md.css ./css/hpre Lect-01.html ./css/hpost >/tmp/Lect-01.html
	mv /tmp/Lect-01.html ./Lect-01.html

run:
	./run.sh

