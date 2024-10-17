# Kavka: Open-Source, Private & Secure Messenger

Kavka is a secure, open-source messaging platform built on GoLang and adhering to Clean Architecture principles. It provides a robust, scalable, and privacy-focused foundation for developing private chat experiences, secure collaboration tools, or confidential messaging apps.

### Built for Privacy

Kavka is more than just a messenger; it's a commitment to your privacy. Every message you send is end-to-end encrypted, ensuring that only you and the intended recipient can read its contents. Our open-source code is transparent, allowing you to verify our security claims. With Kavka, you can communicate with confidence, knowing that your conversations are safe and secure.

### Setup Development Server

The server setup guide is explained in detail here.

 ```bash
 git clone --depth 1 https://github.com/kavkaco/Kavka-Core.git
 ```

2. Install dependencies

 ```bash
 go mod tidy
 ```

3. Edit configuration
```bash
vim ./config/config.yml
```

5. Start services

Lets build and start services with `docker-compose`

 ```bash
 sudo docker-compose up -d mongo redis minio

make dev
 ```


### Contribution!

Kavka welcomes contributions from developers around the world. We are committed to fostering a collaborative and inclusive community where everyone can contribute to the project's growth and development.

We invite you to explore Kavka, join our community, and contribute to building a secure and private communication platform for everyone.
