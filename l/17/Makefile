
all: Lect-17.html 

FR=../01/

Lect-17.html: Lect-17.md
	markdown-cli --input=./Lect-17.md --output=Lect-17.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Lect-17.html ${FR}/css/hpost >/tmp/Lect-17.html
	mv /tmp/Lect-17.html ./Lect-17.html

