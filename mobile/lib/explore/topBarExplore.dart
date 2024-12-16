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
          child: TextTopBar(text: "All", index: 0),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.index = 1;
            onPressed();
          },
          child: TextTopBar(text: "Workflows", index: 1,),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.index = 2;
            onPressed();
          },
          child: TextTopBar(text: "Services", index: 2,),
        ),
      ],
    );
  }
}

class TextTopBar extends StatelessWidget {
  const TextTopBar({
    super.key,
    required this.text,
    this.index = 0,
  });

  final String text;
  final int index;

  @override
  Widget build(BuildContext context) {
    if (index != SelectedCardList.index) {
      return AutoSizeText(
        text,
        style: TextStyle(
            fontSize: 20, fontWeight: FontWeight.w600, color: Colors.black),
        maxLines: 1,
      );
    }
    return AutoSizeText(
        text,
        style: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.w900,
          color: Colors.black,
          decoration: TextDecoration.underline,
          decorationThickness: 2,
        ),
        maxLines: 1,
      );
  }
}
