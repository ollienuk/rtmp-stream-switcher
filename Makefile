APP_NAME=switch
INSTALL_DIR=/opt/$(APP_NAME)
SERVICE_FILE=/etc/systemd/system/$(APP_NAME).service

.PHONY: all install build service start stop status clean uninstall

all: install

install: build copy-files service start

build:
	go build -o $(APP_NAME) main.go

copy-files:
	sudo mkdir -p $(INSTALL_DIR)
	sudo cp $(APP_NAME) config.json $(INSTALL_DIR)/

service:
	echo "[Unit]" | sudo tee $(SERVICE_FILE) > /dev/null
	echo "Description=Stream Service" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "After=network.target" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "[Service]" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Type=simple" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Restart=always" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "RestartSec=3" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "User=root" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Group=root" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "WorkingDirectory=$(INSTALL_DIR)/" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "ExecStart=$(INSTALL_DIR)/$(APP_NAME)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "[Install]" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "WantedBy=multi-user.target" | sudo tee -a $(SERVICE_FILE) > /dev/null

start:
	sudo systemctl daemon-reexec
	sudo systemctl enable $(APP_NAME).service
	sudo systemctl restart $(APP_NAME).service

stop:
	sudo systemctl stop $(APP_NAME).service

status:
	sudo systemctl status $(APP_NAME).service

clean:
	rm -f $(APP_NAME)

uninstall: stop
	sudo systemctl disable $(APP_NAME).service
	sudo rm -f $(SERVICE_FILE)
	sudo rm -rf $(INSTALL_DIR)
