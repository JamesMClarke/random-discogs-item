make:
	conda env create -f ENV.yml
	conda activate discogs
	pip install -e .

remove:
	conda env remove -n discogs
	rm -rf random_discogs_item.egg-info