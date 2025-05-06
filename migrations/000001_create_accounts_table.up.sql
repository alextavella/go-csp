CREATE TABLE IF NOT EXISTS accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 0
);

INSERT INTO accounts (user_id, balance) VALUES ("alextavella", 1000) ON DUPLICATE KEY UPDATE balance=balance;
