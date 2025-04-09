SELECT 
    DATE(created_at) as date, 
    count(distinct playfab_id) as unique_entries,
    count(*) as total_entries,
    promo_id 
FROM promo_entries
WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY DATE(created_at), promo_id
ORDER BY date 