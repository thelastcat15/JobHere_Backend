import json

with open("seeds/parking.json", "r", encoding="utf-8") as f:
    data = json.load(f)

output = []
output.append("""
DO $$
DECLARE
    parking_id uuid;
    zone_id uuid;
    j int;
    k int;
BEGIN
""")

for p in data[:30]:
    name = p["title"].replace("'", "''")
    address = p["carpark_address"].replace("'", "''")
    lat = p["lat"]
    lng = p["lng"]
    tel = p["carpark_tel"]

    output.append(f"""
    INSERT INTO parkings
    (name, type, contact, address, description, coordinate_x, coordinate_y, created_at, updated_at)
    VALUES
    (
        '{name}',
        'PUBLIC',
        '{tel}',
        '{address}',
        'Imported parking',
        {lat},
        {lng},
        NOW(),
        NOW()
    )
    RETURNING id INTO parking_id;

    FOR j IN 1..3 LOOP

        INSERT INTO parking_zones
        (parking_id, name, hour_rate, created_at, updated_at)
        VALUES
        (
            parking_id,
            'Zone ' || j,
            30,
            NOW(),
            NOW()
        )
        RETURNING id INTO zone_id;

        FOR k IN 1..10 LOOP

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
""")

output.append("""
END $$;
""")

# save to txt
with open("seeds/mock_parking_sql.txt", "w", encoding="utf-8") as f:
    f.write("\n".join(output))

print("SQL saved to seeds/mock_parking_sql.txt")