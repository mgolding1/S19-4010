
all: Lect-16.html 

FR=../01/

Lect-16.html: Lect-16.md
	markdown-cli --input=./Lect-16.md --output=Lect-16.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Lect-16.html ${FR}/css/hpost >/tmp/Lect-16.html
	mv /tmp/Lect-16.html ./Lect-16.html


