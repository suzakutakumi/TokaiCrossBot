source ./.env
curl https://api.twitter.com/2/users/by/username/$1 -H "Authorization: Bearer $BearerToken"
