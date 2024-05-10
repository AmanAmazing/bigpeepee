--- Users table -----
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL
);

--- Capsules table -----
CREATE TABLE Capsules (
    capsule_id SERIAL PRIMARY KEY,
    user_id INT,  -- Capsule Owner
    title VARCHAR(255),
    description TEXT,
    open_date DATE,
    is_public BOOLEAN,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);


--- Capsule Content table ---- 
CREATE TYPE ContentType AS ENUM ('text','image','video');
CREATE TABLE Capsule_Content (
    content_id SERIAL PRIMARY KEY,
    capsule_id INT,
    content_type ContentType, 
    content_data Text, -- For text, store as text. For images/videos, store file path/url 
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (capsule_id) REFERENCES Capsules(capsule_id)
);

--- Capsule Access ---- 
CREATE TABLE Capsule_Access(
    access_id SERIAL PRIMARY KEY, 
    user_id INT,
    capsule_id INT,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (capsule_id) REFERENCES Capsules(capsule_id)
);



CREATE TABLE Capsule_Comments (
    comment_id SERIAL PRIMARY KEY, 
    capsule_id INT, 
    user_id INT,
    comment_text TEXT,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (capsule_id) REFERENCES Capsules(capsule_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TYPE NotificationType AS ENUM ('decision_made','comment_added');
CREATE TABLE Notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id INT, 
    capsule_id INT,
    notification_type NotificationType,
    is_read BOOLEAN DEFAULT FALSE,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (capsule_id) REFERENCES Capsules(capsule_id)
);

-- Need application logic to insert notification into this table when capsules open_date becomes the current date 
-- Insert a nottification when a new comment is added to a capsule 

