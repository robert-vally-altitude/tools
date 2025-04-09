SELECT 
    DATE(pe.created_at) as date, 
    count(pe.id) as entries, 
    p.sweepz_entry as sweepz_per_entry, 
    count(pe.id) * p.sweepz_entry as total_sweepz, 
    pe.promo_id, 
    pe.playfab_id
FROM promo_entries pe
LEFT JOIN promos p on pe.promo_id = p.id
WHERE pe.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 14 DAY)
GROUP BY DATE(pe.created_at), pe.promo_id, pe.playfab_id
ORDER BY date, pe.playfab_id 