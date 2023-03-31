# WhatChatGPT

## Whatsapp Chatbot with Golang and Whatsmeow  
This project is an open-source chatbot for Whatsapp, built with Golang and the Whatsmeow library. It uses OpenAI to generate responses, enabling it to have complex conversations. This repo provides you with all the tools required for getting your own chatbot up and running.

## Overview 
The Whatsapp Chatbot with Golang & OpenAI is a project that enables users to create their own conversational AI bot for Whatsapp. It takes advantage of the OpenAI platform to train its own model and use it to generate responses based on user inputs. This project also features an advanced image recognition component that can identify objects in images sent by the user. The Whatsmeow library simplifies the entire process by providing users with all the necessary tools to create powerful chatbots.  

## Features
* Send message as input, receive response as output 
* Send image as input, receive AI generated anime image as output

## Hosting Options 
You can deploy your chatbot on any hosting provider of your choice. For example, if you're using Amazon Web Services (AWS), you can host your chatbot using AWS Lambda. Additionally, you can also host your chatbot using a .service file format. To do this, you will need to add your code to the .service file, configure your environment and set up the appropriate permissions to allow the service to run. You may also need to adjust the configuration settings of your operating system. 

Once all the setup and configurations are complete, you can then run your chatbot using the .service file and start interacting with users. 

We hope that this tutorial helps you get started with creating your own Whatsapp chatbot application. For more information and other tutorials, please check out our website. Thanks!


### Requirements 
- Golang (1.18 or higher) 
- Whatsmeow library 
- A hosting provider for deployment 


### Prerequisites
You'll need access to Golang, the Whatsmeow library, an access token from https://mtlab.meitu.com/ and an API Key from https://platform.openai.com/account/api-keys.

### Setting Up Variables
Add the following variables to your code: 

* `OPENAIKEY` - the API key obtained from OpenAI
* `IMGAIKEY` - the access token obtained from mtlab.meitu.com
* `IMGAISECRET` - the secret associated with the access token obtained from mtlab.meitu.com

### Obtaining Access Token and API Key
To obtain the access token from mtlab.meitu.com, go to the website, log in to your Meitu account and follow their instructions. 

To obtain the API key from OpenAI, you must register for an OpenAI account. Then, you can locate and copy the API key from the Keys & Tokens tab in the account settings page. 


## Installation 
1. Clone the repository: `git clone https://github.com/TAMILVIP007/WhatChatGPT` 
2. Navigate to the project directory: `cd WhatChatGPT` 
3. Install Whatsmeow: `go mod tidy`
4. Build the project: `go build .` 
5. Deploy the bot to your chosen hosting provider 

## Documentation 
For more detailed instructions and documentation, please refer to the README file in the project repo.  

## Hosting with a .service File 
If you want to host this chatbot as a service running in the background, you can create a .service file that contains all the information needed to start the chatbot. This file will be installed in /etc/systemd/system/, where it will be run automatically when the system starts.

You will need to include the following information in the file:
• The path to the executable of your chatbot.
• The user and group that will own the process.
• Any environment variables needed by the chatbot.
• Any other settings or flags that are required.

Once the .service file is created, you can enable it to start up at boot time by running the command ‘sudo systemctl enable [service name]’. You can then start or stop the service with the commands ‘sudo systemctl start [service name]’ and ‘sudo systemctl stop [service name]’.


## Live Demo 
You can try out a live demo of the chatbot [here](https://wa.me/+18739002988).

## Support 
If you encounter any issues when setting up or running this repository, please contact us for support via our support chat at [Bots Realm](https://t.me/mybotsrealm).

## License & Credits 
The code in this repo has been written under the MIT license and all credit for the code should be attributed to the original author. Additionally, credit should be given to [Whatsmeow](https://github.com/tulir/whatsmeow) and OpenAI for their respective technologies used within the application.

