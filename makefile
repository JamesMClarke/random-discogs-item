# Makefile
CONDA_BASE := $(shell bash -c 'source $$HOME/miniconda3/etc/profile.d/conda.sh >/dev/null 2>&1 || true; conda info --base 2>/dev/null')
CONDA := $(CONDA_BASE)/bin/conda

debug:
	@echo "Detected conda: $(CONDA)"

make:
	$(CONDA) env create -f ENV.yml
	$(CONDA) run -n discogs pip install -e .

remove:
	$(CONDA) env remove -n discogs
	rm -rf random_discogs_item.egg-info
