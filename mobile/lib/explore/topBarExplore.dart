import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "explorePageData.dart";

class TopBarExplore extends StatelessWidget {
  TopBarExplore({
    super.key,
    required this.onPressed,
  });

  final Function onPressed;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {
            SelectedCardList.index = 0;
            onPressed();
          },
          child: AutoSizeText(
            "All",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.index = 1;
            onPressed();
          },
          child: AutoSizeText(
            "Workflows",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.index = 2;
            onPressed();
          },
          child: AutoSizeText(
            "Services",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}
