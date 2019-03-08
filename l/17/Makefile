
all: Lect-17.html HW5-note.html

FR=../01/

Lect-17.html: Lect-17.md
	markdown-cli --input=./Lect-17.md --output=Lect-17.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre Lect-17.html ${FR}/css/hpost >/tmp/Lect-17.html
	mv /tmp/Lect-17.html ./Lect-17.html

HW5-note.html: HW5-note.md
	markdown-cli --input=./HW5-note.md --output=HW5-note.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ../../md.css ${FR}/css/hpre HW5-note.html ${FR}/css/hpost >/tmp/HW5-note.html
	mv /tmp/HW5-note.html ./HW5-note.html

