

format:
	js-beautify --replace static/script.js
	go fmt server.go

logs:
	ssh root@simple-web-analytics.com supervisorctl tail -f webstats stderr

deploy:
	ssh root@simple-web-analytics.com 'sh -c "cd webstats && git pull && supervisorctl restart webstats"'

fastdeploy:
	git commit -am - --allow-empty
	git push
	make deploy
