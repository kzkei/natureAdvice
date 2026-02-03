from etl.utils import conn_ops
from pathlib import Path
from psycopg2.extras import execute_values
import logging, json

# set up logging
logger = logging.getLogger(__name__)

# seed 63 nat parks
# initial, undocumented 'docker exec' seed included (5) west parks of park codes: ['GCNP', 'ZNP', 'GNP', 'YNP', 'YosemiteNP']

# use defined json master list
# seed master list with idempotency (initial 5 seeded parks will not fail and will not duplicate)

# load from json
def load_parks():
    logger.info("begin LOADING parks json")

    park_json = Path(__file__).parent / 'parks.json'
    
    if not park_json.exists():
        raise FileNotFoundError(f"parks data file not found: {park_json}")
    
    with open(park_json, 'r') as f:
        parks = json.load(f)
    
    logger.info(f"loaded {len(parks)} parks from {park_json.name}")

    return parks

# seed into db - locations
def seed_parks():
    logger.info("begin SEEDING parks json")

    parks = load_parks()

    # connect to db
    conn, cursor = conn_ops.open()

    try:
        logger.info("attempting seed")

        cursor.execute("SELECT COUNT(*) FROM locations")
        current_count = cursor.fetchone()[0]
        logger.debug(f"{current_count} park(s) currently in db")

        # insert with idempotency (on conflict for existing unique entries)
        insert_query = """
            INSERT INTO locations 
            (name, park_code, latitude, longitude, state, region, timezone, elevation_ft)
            VALUES %s
            ON CONFLICT (park_code) 
            DO UPDATE SET
                name = EXCLUDED.name,
                latitude = EXCLUDED.latitude,
                longitude = EXCLUDED.longitude,
                state = EXCLUDED.state,
                region = EXCLUDED.region,
                timezone = EXCLUDED.timezone,
                elevation_ft = EXCLUDED.elevation_ft,
                updated_at = NOW()
        """

        # convert all parks into tuples
        parks_data = [
            (
                park['name'],
                park['park_code'],
                park['latitude'],
                park['longitude'],
                park['state'],
                park['region'],
                park['timezone'],
                park['elevation_ft']
            )
            for park in parks
        ]

        logger.info("executing insert query")
        execute_values(cursor, insert_query, parks_data)
        conn.commit()
        
        # get new count
        cursor.execute("SELECT COUNT(*) FROM locations")
        after_count = cursor.fetchone()[0]

        logger.debug(f"seeding executed with a result of {after_count} total entries")

    except Exception as e:
        logger.error(f"error seeding json: {e}")
        conn.rollback()
        raise
    finally:
        conn_ops.close(conn=conn, cursor=cursor)


def main():
    logger.info("starting seeding process")

    seed_parks()

    logger.info("seeding complete")


if __name__ == "__main__":
    main()