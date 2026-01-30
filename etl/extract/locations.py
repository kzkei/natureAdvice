import psycopg2
from utils import conn_ops

# EXTRACT locations from locations table (seeded and or added via API CREATE (*))
def extract_location_metadata():
    # TODO logging, error handling

    # open psycop conn
    conn, cursor = conn_ops.open()

    # execute fetch SQL into locations
    cursor.execute("SELECT * FROM locations;")
    locations = cursor.fetchall()

    # close
    conn_ops.close(conn=conn, cursor=cursor)

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