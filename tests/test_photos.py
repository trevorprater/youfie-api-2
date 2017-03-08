import unittest
import requests
import json
from utils import utils


class TestPhoto(unittest.TestCase):
    def setUp(self):
        self.photo = {
            'format': 'jpg',
            'storage_url': 'http://i.imgur.com/MsCzsxP.jpg',
            'latitude': 40.744381,
            'longitude': -73.987333,
            'width': 2448,
            'height': 3264,
        }
        utils.create_user('trevor', 'trevor@youfie.io', 'venice')
        r, self.session = utils.login('trevor', 'venice')

    def tearDown(self):
        utils.delete_user_if_exists('trevor', 'venice')

    def test_create_photo(self):
        r = utils.create_photo('trevor', self.photo, self.session)
        photo = json.loads(r.content)
        self.assertEqual(r.status_code, 201)
        self.assertEqual(photo['format'], 'jpg')
        self.assertEqual(photo['storage_url'],
                         'http://i.imgur.com/MsCzsxP.jpg')
        self.assertEqual(photo['latitude'], 40.744381)
        self.assertEqual(photo['longitude'], -73.987333)
        self.assertEqual(photo['width'], 2448)
        self.assertEqual(photo['height'], 3264)
        self.assertEqual(photo['processed'], False)
        utils.delete_photo('trevor', photo['id'], self.session)

    def test_view_photo(self):
        r = utils.create_photo('trevor', self.photo, self.session)
        photo = json.loads(r.content)
        r = utils.get_photo('trevor', photo['id'], self.session)
        get_photo = json.loads(r.content)
        self.assertEqual(photo['id'], get_photo['id'])
        self.assertEqual(self.photo['storage_url'], photo['storage_url'])
        self.assertEqual(self.photo['format'], get_photo['format'])
        self.assertEqual(self.photo['latitude'], photo['latitude'])
        self.assertEqual(self.photo['longitude'], photo['longitude'])
        self.assertEqual(self.photo['width'], photo['width'])
        self.assertEqual(self.photo['height'], photo['height'])
        utils.delete_photo('trevor', get_photo['id'], self.session)

    def test_delete_photo(self):
        r = utils.create_photo('trevor', self.photo, self.session)
        photo = json.loads(r.content)
        r = utils.get_photo('trevor', photo['id'], self.session)
        self.assertEqual(r.status_code, 200)
        r = utils.delete_photo('trevor', photo['id'], self.session)
        self.assertEqual(r.status_code, 201)
        r = utils.get_photo('trevor', photo['id'], self.session)
        self.assertEqual(r.status_code, 404)
