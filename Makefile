cleanup:
	k3d cluster delete dev
	k3d cluster create dev

vmail:
	docker build -t vmail:latest -f vmail.Dockerfile .
	k3d image import vmail:latest -c dev
	kubectl apply -f deploy/vmail-pvc.yaml --force
	kubectl apply -f deploy/vmail.yaml --force

vmail-test:
	./test-vmail.sh

vapi:
	docker build -t vapi:latest -f vapi.Dockerfile .
	k3d image import vapi:latest -c dev
	kubectl apply -f deploy/vapi-pvc.yaml --force
	kubectl apply -f deploy/vapi.yaml --force

postfix:
	kubectl apply -f deploy/postfix.yaml --force

e2e:
	# echo -e "EHLO test.local\r\nMAIL FROM: sender@example.com\r\nRCPT TO: test@vmail.today\r\nDATA\r\nSubject: Test Email\r\n\r\nThis is a test email sent to Postfix.\r\n.\r\nQUIT\r\n" | nc localhost 2525
	./e2e.sh

www:
	docker build -t www:latest -f www.Dockerfile .
	k3d image import www:latest -c dev
	
local:
	kubectl port-forward svc/postfix 2525:25	
	kubectl port-forward deployment/www-deployment 8090:80
