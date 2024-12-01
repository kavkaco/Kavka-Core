require "securerandom"

token_length = 64

new_token = SecureRandom.base64(token_length)

puts new_token
