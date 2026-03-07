DO $$
DECLARE
    parking_id uuid;
    zone_id uuid;
    i int;
    j int;
    k int;
BEGIN

FOR i IN 1..10 LOOP

    INSERT INTO parkings
    (name, type, contact, address, description, coordinate_x, coordinate_y, created_at, updated_at)
    VALUES
    (
        'Parking ' || i,
        'PUBLIC',
        '0123456789',
        'Mock Address ' || i,
        'Mock parking description',
        13.7563 + random()/100,
        100.5018 + random()/100,
        NOW(),
        NOW()
    )
    RETURNING id INTO parking_id;

    -- create zones
    FOR j IN 1..3 LOOP

        INSERT INTO parking_zones
        (parking_id, name, hour_rate, created_at, updated_at)
        VALUES
        (
            parking_id,
            'Zone ' || j,
            (random()*50 + 20)::numeric(10,2),
            NOW(),
            NOW()
        )
        RETURNING id INTO zone_id;

        -- create slots
        FOR k IN 1..20 LOOP

            INSERT INTO parking_slots
            (zone_id, name, status, created_at, updated_at)
            VALUES
            (
                zone_id,
                'SLOT-' || k,
                'available',
                NOW(),
                NOW()
            );

        END LOOP;

    END LOOP;

END LOOP;

END $$;