curl -H "login:john" -H "password:doe" \
    -X POST -d '{"greeting":"foo"}' 'http://dragon:7778/1/ping'
