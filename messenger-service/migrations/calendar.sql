CREATE TABLE meetings (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    owner_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    reminder_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (chat_id) REFERENCES chat.chats(id) ON DELETE CASCADE,
    FOREIGN KEY (owner_id) REFERENCES chat.users(id)
);

CREATE TABLE meeting_participants (
    meeting_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (meeting_id, user_id),
    FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES chat.users(id)
);