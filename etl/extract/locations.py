import psycopg2, logging
from utils import conn_ops

# setup logging
logger = logging.getLogger(__name__)

# EXTRACT locations from locations table (seeded and or added via API CREATE (*))
# returns extracted locations
def extract_location_data():
    logger.info("Starting location extraction")

    # open psycop conn
    conn, cursor = conn_ops.open()

    # execute fetch SQL into locations
    cursor.execute("SELECT * FROM locations;")
    locations = cursor.fetchall()

    logger.debug(f"locations extracted: {locations}")

    # close
    conn_ops.close(conn=conn, cursor=cursor)

    logger.info("location extraction complete")

    # return converted tuples in list of dict
    return [
        {
            'id': row[0],
            'name': row[1],
            'park_code': row[2],
            'latitude': row[3],
            'longitude': row[4],
            'state': row[5],
            'region': row[6],
            'elevation_ft': row[7],
            'timezone': row[10]
        }
        for row in locations
    ]