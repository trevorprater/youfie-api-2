import unittest
import requests
import json
from utils import utils


class TestCreatePhoto(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_create_photo_passes(self):
        photo = {
            'format': 'jpg',
            'storage_url': 'http://i.imgur.com/MsCzsxP.jpg',
            'latitude': 40.744381,
            'longitude': -73.987333,
            'width': 2448,
            'height': 3264,
        }
        utils.create_user('trevor', 'trevor@youfie.io', 'venice')
        r, session = utils.login('trevor', 'venice')
        r = utils.create_photo('trevor', photo, session)
        self.assertEqual(r.status_code, 201)


