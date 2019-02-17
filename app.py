from flask import Flask, request, Response
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
import json
import hashlib
import base64

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///yuefan.sqlite3'

db = SQLAlchemy(app)
migrate = Migrate(app, db)

# Models
restaurant_restrict = db.Table('restaurant_restrict',
                    db.Column('user_id', db.Integer, db.ForeignKey('user.id')),
                    db.Column('restaurant_id', db.Integer, db.ForeignKey('restaurant.id'))
)


class Restaurant(db.Model):
    id = db.Column('id', db.Integer, primary_key=True)
    name = db.Column(db.String(10))
    status = db.Column(db.Boolean())

    def __init__(self, name):
        self.name = name
        self.status = False


class User(db.Model):
    id = db.Column('id', db.Integer, primary_key=True)
    username = db.Column(db.String(20))
    password = db.Column(db.String(32))
    restaurants = db.relationship('Restaurant',
                                  secondary=restaurant_restrict,
                                  backref=db.backref('user'))

    def __init__(self, username, password):
        self.username = username
        self.password = password


# Router
@app.route('/api/add_restaurant', methods=['POST'])
def add_restaurant():
    restaurant = request.json.get('restaurant')
    username = request.json.get('username')
    user = User.query.filter_by(username=username).first()
    exist_restaurant = Restaurant.query.filter_by(name=restaurant).first()
    add_restaurant = Restaurant(restaurant)
    if exist_restaurant:
        add_restaurant = exist_restaurant
    user.restaurants.append(add_restaurant)
    db.session.commit()
    return Response(json.dumps({
        'status': True,
        'restaurant': restaurant
    }), mimetype='application/json')


@app.route('/api/get_restaurants', methods=['POST'])
def get_restaurants():
    restaurant_set = Restaurant.query.filter_by(status=True).all()
    restaurants = []
    user_restaurants = []
    username = request.json.get('username')
    user = User.query.filter_by(username=username).first()
    for restaurant in restaurant_set:
        restaurants.append(restaurant.name)
    if username:
        for restaurant in user.restaurants:
            if restaurant.name not in restaurants:
                restaurants.append(restaurant.name)
            user_restaurants.append(restaurant.name)
    return Response(json.dumps({
        'restaurants': restaurants,
        'user_restaurants': user_restaurants
    }), mimetype="application/json")


@app.route('/api/delete_restaurants', methods=['POST'])
def delete_restaurants():
    restaurants = request.json.get('restaurants')
    username = request.json.get('username')
    user = User.query.filter_by(username=username).first()
    remain_restaurants = []
    for restaurant in restaurants:
        deleted_restaurant = Restaurant.query.filter_by(name=restaurant['name']).first()
        user.restaurants.remove(deleted_restaurant)
        exist_user = User.query.filter(User.restaurants.contains(deleted_restaurant))
        if exist_user.count() == 0 and not deleted_restaurant.status:
            db.session.delete(deleted_restaurant)
    db.session.commit()
    for restaurant in user.restaurants:
        remain_restaurants.append({'name': restaurant.name})
    return Response(json.dumps({
        'status': True,
        'msg': 'Deletion has done successfully!',
        'restaurants': remain_restaurants
    }))


@app.route('/api/add_user', methods=['POST'])
def add_user():
    username = request.json.get('username')
    password = request.json.get('password')
    password = str(base64.b64decode(password), 'utf-8')
    people = User.query.filter_by(username=username).all()
    if not people:
        hl = hashlib.md5()
        hl.update(password.encode(encoding='utf-8'))
        db.session.add(User(username, hl.hexdigest()))
        db.session.commit()
        return Response(json.dumps({
            'status': True,
            'msg': 'Successfully registered! The page will redirect to login in 3s.'
        }))
    return Response(json.dumps({
        'status': False,
        'msg': 'The username has registered'
    }), mimetype='application/json')


@app.route('/api/login', methods=['POST'])
def login():
    username = request.json.get('username')
    password = request.json.get('password')
    password = str(base64.b64decode(password), 'utf-8')
    people = User.query.filter_by(username=username).first()
    if not people:
        return Response(json.dumps({
            'status': False,
            'msg': 'User not exist'
        }))
    hl = hashlib.md5()
    hl.update(password.encode(encoding='utf-8'))
    if hl.hexdigest() != people.password:
        return Response(json.dumps({
            'status': False,
            'msg': 'Password not match'
        }))
    return Response(json.dumps({
        'status': True,
        'msg': 'Login Success'
    }))


# Main Function
if __name__ == '__main__':
    app.run()
