from etl.utils import conn_ops

# seed 63 nat parks
# initial, undocumented seed included (5) west parks of park codes: ['GCNP', 'ZNP', 'GNP', 'YNP', 'YosemiteNP']

# use defined master list json
# seed master list with idempotency (initial 5 seeded parks will not fail and will not duplicate)

# load from json
def load_parks():
    # TODO logging, error handling

    return

# seed into db - locations
def seed_parks():
    # TODO logging, error handling

    parks = load_parks()

    # connect
    conn, cursor = conn_ops.open()

    # try: 

    return 

    


