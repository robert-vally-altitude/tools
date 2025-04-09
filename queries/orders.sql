SELECT 
    playfab_id, 
    pspReference, 
    country_code, 
    order_number, 
    amount, 
    game_id, 
    playfab_item_id, 
    display_name, 
    created_at, 
    updated_at
FROM orders 
WHERE status = 'paid' 