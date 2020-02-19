for i in {0..500}
do
	curl -X POST -H "X-Parse-Application-Id: appId" -H "Content-Type: application/json" -d '{"score":1337}' http://localhost:1337/parse/classes/test
done
