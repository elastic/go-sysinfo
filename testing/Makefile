# Control ansible's verbosity by setting this with -v or -vvvvv (more verbose).
ANSIBLE_VERBOSE?=
# Exluding localhost, as localhost is assumed to be Mac OS X with ssh enabled.
ANSIBLE_LIMIT?=all:!localhost
# Extra flags to pass to Ansible (e.g. --skip-tags filebeat).
ANSIBLE_EXTRA_FLAGS?=--tags run_tests

.PHONY:ansible
ansible: ve
	ANSIBLE_HOST_KEY_CHECKING=false \
	ANSIBLE_SSH_ARGS='-o UserKnownHostsFile=/dev/null -o IdentitiesOnly=yes -o ControlMaster=auto -o ControlPersist=60s' \
	time ve/bin/ansible-playbook \
	${ANSIBLE_VERBOSE} \
	-i .vagrant/provisioners/ansible/inventory/vagrant_ansible_inventory \
	--limit "${ANSIBLE_LIMIT}" \
	${ANSIBLE_EXTRA_FLAGS} \
	ansible/main.yml

# Create a virtualenv to run Ansible.
ve: ve/bin/activate
ve/bin/activate: requirements.txt
	@test -d ve || virtualenv ve
	@ve/bin/pip install -Ur requirements.txt
	@touch ve/bin/activate

.PHONY:setup
setup:
	vagrant up

# This destroys all vagrant machines and removes the vagrant related data
clean:
	-vagrant destroy -f
	-rm -r .vagrant
	-rm -f ssh_config
