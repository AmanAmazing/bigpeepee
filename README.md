# Big Peepee project 
Building this project in my spare time so don't judge. I ain't no expert :) 

## Tech stack 
I call it the big peepee stack - Golang, templ, htmx :) 
No slowy javascript (maybe later and htmx uses it but it is what it is)
Maybe some react for dynamic stuff - Big maybe but I don't think I will need it. 


## Docker instructions
remember to delete any existing postgres-data folder that may have accidentally been left over. I don't do database migrations, I just rebuild it lol
`docker compose build` `docker compose up` 

### Running it in the background and displaying logs
`docker-compose up -d`
`docker-compose logs -f go-app` // the 'go-app' is replaced by whatever application logs you want. I will dockerize this go-app once I am done coding it. 


### PG admin 
This requires an initial setup stage  
Even with these environment variables, pgAdmin does not support automatically creating server connections. You will need to manually configure the connection when you first access pgAdmin:
    1)Open pgAdmin by navigating to http://localhost:5050 in your browser.
    2)Log in using the PGADMIN_DEFAULT_EMAIL and PGADMIN_DEFAULT_PASSWORD.
    3)Right-click on 'Servers' in the left-hand browser pane and select 'Create' -> 'Server'.
    4)In the 'Create - Server' dialog:
        - General tab: Enter a name for your server (e.g., "Postgres Server").
        - Connection tab:
            - Hostname/address: Enter postgres (the name of the service in docker-compose.yml).
            - Port: 5432 (or your custom PostgreSQL port).
            - Maintenance database: Enter the database name (e.g., "posystem").
            - Username: Enter the PostgreSQL username (e.g., "aman").
            - Password: Enter the PostgreSQL password (e.g., "admin").
    5) Click 'Save'.

## Initial system design This can be seen in the file ./init-system-design.svg the ./postgres-init/init.sql shows the database structure. I ain't no fucking database expert so don't come for me! 
I tried to make the system semi flexible :) just realised the smiley face looks bloody wonky in the font on vim (use vim btw). 
TODO! Need to add a total price section in the purchase_orders table
TODO! Need to add notification section to the database init file 
