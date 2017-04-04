function run_script {
  #Facebook
  echo "Doing Facebook..."
  if [ -n "$facebook" ]; then
    apiclient -api facebook -fb-page $facebook
  fi

  #Twitter Search
  echo "Doing Twitter Search..."
  if [ -n "$twitter_search" ]; then
    apiclient -api twitter -twitter-search=true -twitter-query "\"$twitter_search\""
  fi

  #Twitter timeline
  echo "Doing Twitter Timeline..."
  if [ -n "$twitter_timeline" ]; then
    apiclient -api twitter -twitter-timeline=true -twitter-screenname "$twitter_timeline"
  fi

  #Event Name - Trends, News, Google Search, Google First Page
  echo "Doing Event Name..."
  if [ -n "$event_name" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_name\""
    #apiclient -api google -google-search-query "\"$event_name\""
    #apiclient -api first-page -google-search-query "\"$event_name\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_name\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_name\""
    sleep 5
  fi

  #Event Venue - Trends, News, Google Search, Google First Page
  echo "Doing Event Venue..."
  if [ -n "$event_venue" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_venue\""
    #apiclient -api google -google-search-query "\"$event_venue\""
    #apiclient -api first-page -google-search-query "\"$event_venue\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_venue\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_venue\""
    sleep 5
  fi

  #Event City - Trends, News, Google Search, Google First Page
  echo "Doing Event City..."
  if [ -n "$event_city" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_city\""
    #apiclient -api google -google-search-query "\"$event_city\""
    #apiclient -api first-page -google-search-query "\"$event_city\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_city\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_city\""
    sleep 5
  fi

  #Event Name+Year - Trends, News, Google Search, Google First Page
  echo "Doing Event Name+Year..."
  if [ -n "$event_name" ] && [ -n "$event_year" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_name\"+\"$event_year\""
    #apiclient -api google -google-search-query "\"$event_name\"+\"$event_year\""
    #apiclient -api first-page -google-search-query "\"$event_name\"+\"$event_year\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_name\"+\"$event_year\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_name\"+\"$event_year\""
    sleep 5
  fi

  #Event Name+Venue - Trends, News, Google Search, Google First Page
  echo "Doing Event Name+Venue..."
  if [ -n "$event_name" ] && [ -n "$event_venue" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_name\"+\"$event_venue\""
    #apiclient -api google -google-search-query "\"$event_name\"+\"$event_venue\""
    #apiclient -api first-page -google-search-query "\"$event_name\"+\"$event_venue\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_name\"+\"$event_venue\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_name\"+\"$event_venue\""
    sleep 5
  fi

  #Event Name+Country - Trends, News, Google Search, Google First Page
  echo "Doing Event Name+Country..."
  if [ -n "$event_name" ] && [ -n "$event_venue" ]; then
	echo "Bing"
    apiclient -api bing -bing-search-query "\"$event_name\"+\"$event_country\""
    #apiclient -api google -google-search-query "\"$event_name\"+\"$event_country\""
    #apiclient -api first-page -google-search-query "\"$event_name\"+\"$event_country\""
	echo "News"
    apiclient -api news -news-search-query "\"$event_name\"+\"$event_country\""
	sleep 5
	echo "Trends"
    apiclient -api trends -trends-search-query "\"$event_name\"+\"$event_country\""
    sleep 5
  fi
}

location=${0:2}
location2=${location: -10}
len=`expr length $location`
len2=`expr length $location2`
diff=$(( $len-$len2 ))
FOLDERS="$PWD/${location:0:$diff}lisbon"

for folder in $FOLDERS
do
  for file in $folder/*.sh
  do
    echo $file
    source $file
    run_script
  done
done
