from flask import Flask, request, url_for, redirect, Response
from flask_sqlalchemy import SQLAlchemy
import json

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///yuefan.sqlite3'

db = SQLAlchemy(app)


class Restaurants(db.Model):
    id = db.Column('id', db.Integer, primary_key=True)
    name = db.Column(db.String(10))
    status = db.Column(db.Boolean())

    def __init__(self, name):
        self.name = name
        self.status = False


@app.route('/')
def root():
    rsp = {
        'name': 'ecwu',
        'status': True
    }
    return Response(json.dumps(rsp), mimetype='application/json')


if __name__ == '__main__':
    db.create_all()
    app.run(debug=True)
