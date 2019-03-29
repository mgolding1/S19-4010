
all: Lect-21.html t1.html

FR=../01/

Lect-21.html: Lect-21.md
	markdown-cli --input=./Lect-21.md --output=Lect-21.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Lect-21.html ${FR}/css/hpost >/tmp/Lect-21.html
	mv /tmp/Lect-21.html ./Lect-21.html

t1.html: t1.md
	markdown-cli --input=./t1.md --output=t1.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre t1.html ${FR}/css/hpost >/tmp/t1.html
	mv /tmp/t1.html ./t1.html



