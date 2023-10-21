#!/bin/bash

FILE="./scratchpad"

# could probably make this so it loads al of the text into a var, but let us be real
# we dont have to think about optimisations like that right now
TITLE=$(head -n 1 $FILE)
DATE=$(head -n 2 $FILE | tail -n 1)
AUTHOR=$(head -n 3 $FILE | tail -n 1)
DESC=$(head -n 4 $FILE | tail -n 1)
BLOGPATH=$(head -n 5 $FILE | tail -n 1)
TOPICS=$(head -n 6 $FILE | tail -n 1)
NOTES=$(head -n 7 $FILE | tail -n 1)

SINGLEQUOTE="'"
DOUBLESINGLE="''"
sed -in "s/$SINGLEQUOTE/$DOUBLESINGLE/g" $FILE
CONTENT=$(sed -n '8,$p' $FILE)

echo $TITLE, $AUTHOR, $DATE, $DESC, $BLOGPATH, $TOPICS, $NOTES

echo $CONTENT

FINAL_STATEMENT="INSERT INTO blogs (blogtitle, blogdate, blogauthor, blogdescription, blogpathname, blogtopics, blognotes, blogcontent) VALUES ('$TITLE', '$DATE', '$AUTHOR', '$DESC', '$BLOGPATH', '$TOPICS', '$NOTES', '$CONTENT');"

echo ""
echo $FINAL_STATEMENT
