SELECT 
    DATE(ev.email_verified_at), 
    ev.email, 
    ev.playfab_id, 
    pe.promo_id
FROM email_verifications ev
INNER JOIN promo_entries pe ON pe.playfab_id = ev.playfab_id
WHERE DATE(ev.email_verified_at) = ?{date} 