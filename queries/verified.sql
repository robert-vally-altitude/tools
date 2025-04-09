SELECT 
    DATE(a.created_at) as created_at,
    COUNT(*) as num_created,
    SUM(CASE WHEN email_verified_at IS NOT NULL THEN 1 ELSE 0 END) as num_verified
FROM (
    SELECT playfab_id, DATE(created_at) as created_at, email_verified_at
    FROM email_verifications
    WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
    GROUP BY playfab_id, DATE(created_at)
) as a
GROUP BY DATE(created_at)
ORDER BY created_at 