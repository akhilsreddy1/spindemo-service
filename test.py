
class Imdb(db.EmbeddedDocument):
    imdb_id = db.StringField()
    rating = db.DecimalField()
    votes = db.IntField()        
    
class Director(db.DynamicDocument):
    pass

class Cast(db.DynamicEmbeddedDocument):
    Name = db.StringField()
    Role = db.StringField()


class Movie(db.Document):
    title = db.StringField(required=True)
    year = db.IntField()
    rated = db.StringField()
    director = db.ReferenceField(Director)
    cast = db.EmbeddedDocumentField(Cast)  # This has all the properties of DynamicDocument and EmbeddedDocument
    poster = db.FileField()
    imdb = db.EmbeddedDocumentField(Imdb)


 
@app.route('/movies/', methods=["POST"])
def add_movie():
    body = request.get_json()
    director_id = body.pop('director', None)
    if director_id and ObjectId.is_valid(director_id):
        director = Director.objects.get(id=ObjectId(director_id))
        body['director'] = director


    movie = Movie(**body).save()
    return jsonify(movie), 201

@app.route('/movies-embed/', methods=["POST"])
def add_movie_embed():
    # Created Imdb object
    imdb = Imdb(imdb_id="12340mov", rating=4.2, votes=7.9)
    body = request.get_json()
    director_id = body.pop('director', None)
    # Add object to movie and save
    movie = Movie(imdb=imdb, **body).save()
    return jsonify(movie), 201

@app.route('/movies/', methods=["GET"])
def get_movies():
    movies = Movie.objects().to_json()
    return Response(movies, mimetype="application/json", status=200)

@app.route('/director/', methods=['POST'])
def add_dir():
    body = request.get_json()
    director = Director(**body).save()
    return jsonify(director), 201

@app.route('/movies/<id>', methods=["PATCH"])
def update_movie(id):
    body = request.get_json()
    director_id = body.pop('director', None)
    if director_id and ObjectId.is_valid(director_id):
        director = Director.objects.get(id=ObjectId(director_id))
        body['director'] = director
    movie = Movie.objects.get(id=id)
    print(body)
    movie.update(**body)
    return jsonify(movie), 200

app.run(port=5000, debug=True)


