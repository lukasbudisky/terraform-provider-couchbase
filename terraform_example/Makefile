clean:
	rm -rf terraform.tfstate crash.log .terraform.lock.hcl .terraform .terraform.lock terraform.tfstate.backup terraform.tfstate.*.backup terraform.backup

run:
	terraform init
	terraform apply -auto-approve

destroy:
	terraform destroy -auto-approve
	clean
