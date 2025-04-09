SELECT 
    playfab_id, 
    email, 
    email_verified_at, 
    created_at
FROM email_verifications
WHERE email = ?{email} 