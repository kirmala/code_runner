Update Version ansible-playbook site.yml -t update -e "code_runner_sandbox_tag=2.0.0"

Apply Migrations	ansible-playbook site.yml -t migrate

Restart App	ansible-playbook site.yml -t manage -e "app_state=restarted"

Stop App	ansible-playbook site.yml -t manage -e "app_state=stopped"

Full Setup	ansible-playbook site.yml