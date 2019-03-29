
all: Lect-22.html t1.html

FR=../01/

Lect-22.html: Lect-22.md
	markdown-cli --input=./Lect-22.md --output=Lect-22.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Lect-22.html ${FR}/css/hpost >/tmp/Lect-22.html
	mv /tmp/Lect-22.html ./Lect-22.html

t1.html: t1.md
	markdown-cli --input=./t1.md --output=t1.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre t1.html ${FR}/css/hpost >/tmp/t1.html
	mv /tmp/t1.html ./t1.html



